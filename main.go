package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/Kirisakiii/neko-micro-blog-backend/configs"
	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/controllers"
	"github.com/Kirisakiii/neko-micro-blog-backend/crons"
	"github.com/Kirisakiii/neko-micro-blog-backend/loggers"
	"github.com/Kirisakiii/neko-micro-blog-backend/middlewares"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	search "github.com/Kirisakiii/neko-micro-blog-backend/proto"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
)

var (
	logger              *logrus.Logger
	cfg                 *configs.Config
	db                  *gorm.DB
	redisClient         *redis.Client
	mongoClient         *mongo.Client
	searchSeviceConn    *grpc.ClientConn
	searchServiceClient search.SearchEngineClient
	storeFactory        *stores.Factory
	controllerFactory   *controllers.Factory
	middlewareFactory   *middlewares.Factory
)

func init() {
	// logger
	logger = loggers.NewLogger()
	logger.Infoln("正在执行程序初始化...")

	var err error

	// 加载配置文件
	cfg, err = configs.NewConfig()
	if err != nil {
		logger.Panicln(err.Error())
	}

	// 设置日志等级
	var (
		logLevel logrus.Level
		logMode  gormLogger.LogLevel
	)
	switch cfg.Env.Type {
	case "development":
		logLevel = logrus.DebugLevel
		logMode = gormLogger.Error
	case "production":
		logLevel = logrus.InfoLevel
		logMode = gormLogger.Silent
	default:
		logLevel = logrus.InfoLevel
		logMode = gormLogger.Silent
	}

	// 设置logrus日志等级
	logger.SetLevel(logLevel)
	logger.Debugln("日志记录等级设定为:", strings.ToUpper(logLevel.String()))

	// 连接到 pgsql 数据库
	logger.Debugln("尝试连接至数据库...")
	db, err = gorm.Open(
		postgres.Open(fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		)),
		&gorm.Config{
			Logger: gormLogger.Default.LogMode(logMode),
		},
	)
	if err != nil {
		logger.Panicln(err.Error())
	}
	logger.Debugln("数据库连接成功")

	// 迁移模型
	logger.Debugln("正在迁移数据表模型...")
	err = models.Migrate(db)
	if err != nil {
		logger.Panicln("迁移数据库模型失败：", err.Error())
	}

	// 建立 Redis 连接
	logger.Debugln("正在连接至 Redis...")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	_, err = redisClient.Ping(context.TODO()).Result()
	if err != nil {
		logger.Panicln("连接至 Redis 失败：", err.Error())
	}
	logger.Debugln("Redis 连接成功")

	// 建立 MongoDB 连接
	logger.Debugln("正在连接至 MongoDB...")
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		logger.Panicln("连接至 MongoDB 失败：", err.Error())
	}
	err = mongoClient.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		logger.Panicln("连接至 MongoDB 失败：", err.Error())
	}
	logger.Debugln("MongoDB 连接成功")

	// 建立搜索服务 gRPC 连接
	searchSeviceConn, err = grpc.Dial(fmt.Sprintf("%s:%d", cfg.SearchService.Host, cfg.SearchService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Panicln("连接至搜索服务失败：", err.Error())
	}
	searchServiceClient = search.NewSearchEngineClient(searchSeviceConn)

	// 建立数据访问层工厂
	storeFactory = stores.NewFactory(db, redisClient, mongoClient, searchServiceClient)

	// 建立控制器层工厂
	controllerFactory = controllers.NewFactory(
		services.NewFactory(storeFactory),
	)

	// 建立中间件工厂
	middlewareFactory = middlewares.NewFactory(storeFactory)
}

func main() {
	// 初始化定时任务
	crons.InitJobs(logger, db, redisClient)

	// 创建 fiber 实例
	var fiberConfig fiber.Config
	// 如果是生产环境，则开启 Prefork
	if cfg.Env.Type == "production" {
		fiberConfig = fiber.Config{
			Prefork: true,
		}
	}
	fiberConfig.BodyLimit = consts.REQUEST_BODY_LIMIT
	app := fiber.New(fiberConfig)

	// 设置中间件
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[${time}][${latency}][${status}][${method}] ${path}\n",
	}))
	app.Use(compress.New(compress.Config{
		Level: cfg.Compress.Level,
	}))

	// Auth 中间件
	authMiddleware := middlewareFactory.NewTokenAuthMiddleware()

	// 静态资源路由
	resource := app.Group("/resources")
	// 头像资源路由
	resource.Static("/avatar", consts.AVATAR_IMAGE_PATH, fiber.Static{
		Compress: true,
	})
	// 博文图片资源路由
	resource.Static("/image", consts.POST_IMAGE_PATH, fiber.Static{
		Compress: true,
	})

	// api 路由
	api := app.Group("/api")

	// // Token 路由
	// tokenController := controllerFactory.NewTokenController()
	// token := api.Group("/token")
	// token.Get("/check", tokenController.NewCheckTokenHandler())      // 检查令牌可用性
	// token.Post("/refresh", tokenController.NewRefreshTokenHandler()) // 刷新令牌

	// User 路由
	userController := controllerFactory.NewUserController()
	user := api.Group("/user")
	user.Get("/profile", userController.NewProfileHandler())                                             // 查询用户信息
	user.Post("/register", userController.NewRegisterHandler())                                          // 用户注册
	user.Post("/login", userController.NewLoginHandler())                                                // 用户登录
	user.Post("/upload-avatar", authMiddleware.NewMiddleware(), userController.NewUploadAvatarHandler()) // 上传头像
	user.Post("/update-psw", userController.NewUpdatePasswordHandler())                                  // 修改密码
	user.Post("/edit", authMiddleware.NewMiddleware(), userController.NewUpdateProfileHandler())         // 修改用户资料

	// Post 路由
	postController := controllerFactory.NewPostController(searchServiceClient)
	post := api.Group("/post")
	post.Get("/list", postController.NewPostListHandler(storeFactory.NewUserStore()))                                 // 获取文章列表
	post.Get("/user-status", authMiddleware.NewMiddleware(), postController.NewPostUserStatusHandler())               // 获取用户文章状态
	post.Post("/new", authMiddleware.NewMiddleware(), postController.NewCreatePostHandler())                          // 创建文章
	post.Post("/upload/img/file", authMiddleware.NewMiddleware(), postController.NewUploadPostImageFromFileHandler()) // 上传博文图片
	post.Post("/upload/img/url", authMiddleware.NewMiddleware(), postController.NewUploadPostImageFromURLHandler())   // 从 URL 上传博文图片
	post.Post("/like", authMiddleware.NewMiddleware(), postController.NewLikePostHandler())                           // 点赞文章
	post.Post("/cancel-like", authMiddleware.NewMiddleware(), postController.NewCancelLikePostHandler())              // 取消点赞文章
	post.Post("/favourite", authMiddleware.NewMiddleware(), postController.NewFavouritePostHandler())                 // 收藏文章
	post.Post("/cancel-favourite", authMiddleware.NewMiddleware(), postController.NewCancelFavouritePostHandler())    // 取消收藏文章
	post.Get("/:post", postController.NewPostDetailHandler())                                                         // 获取文章信息
	post.Delete("/:post", authMiddleware.NewMiddleware(), postController.NewDeletePostHandler())                      // 删除文章

	// Comment 路由
	commentController := controllerFactory.NewCommentController()
	comment := api.Group("/comment")
	comment.Get("/list", commentController.NewCommentListHandler())                                                     // 获取评论列表
	comment.Get("/detail", commentController.NewCommentDetailHandler())                                                 // 获取评论详情信息
	comment.Get("/user-status", authMiddleware.NewMiddleware(), commentController.NewCommentUserStatusHandler())        // 获取用户评论状态
	comment.Post("/edit", authMiddleware.NewMiddleware(), commentController.NewUpdateCommentHandler())                  // 修改评论
	comment.Post("/delete", authMiddleware.NewMiddleware(), commentController.DeleteCommentHandler())                   // 删除评论
	comment.Post("/like", authMiddleware.NewMiddleware(), commentController.NewLikeCommentHandler())                    // 点赞评论
	comment.Post("/cancel-like", authMiddleware.NewMiddleware(), commentController.NewCancelLikeCommentHandler())       // 取消点赞评论
	comment.Post("/dislike", authMiddleware.NewMiddleware(), commentController.NewDislikeCommentHandler())              // 踩评论
	comment.Post("/cancel-dislike", authMiddleware.NewMiddleware(), commentController.NewCancelDislikeCommentHandler()) // 取消踩评论
	comment.Post("/new", authMiddleware.NewMiddleware(), commentController.NewCreateCommentHandler(
		storeFactory.NewPostStore(),
		storeFactory.NewUserStore(),
	)) // 创建评论

	// Reply 路由
	replyController := controllerFactory.NewReplyController()
	reply := api.Group("/reply")
	reply.Get("/list", replyController.NewGetReplyListHandler())     // 获取回复列表
	reply.Get("/detail", replyController.NewGetReplyDetailHandler()) // 获取回复详情信息
	reply.Post("/new", authMiddleware.NewMiddleware(), replyController.NewCreateReplyHandler(
		storeFactory.NewCommentStore(),
		storeFactory.NewUserStore()),
	) // 创建回复
	reply.Post("/edit", authMiddleware.NewMiddleware(), replyController.NewUpdateReplyHandler()) // 修改回复
	reply.Post("/delete", authMiddleware.NewMiddleware(), replyController.DeleteReplyHandler())  // 删除回复

	// Search 路由
	searchController := controllerFactory.NewSearchController(searchServiceClient)
	search := api.Group("/search")
	search.Get("/post", searchController.NewSearchPostHandler()) // 搜索文章

	// 启动服务器
	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", cfg.Database.Host, cfg.Server.Port)))
}
