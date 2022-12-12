package main

import (
	"fmt"
	"net/http"
	oapi "sample-go-echo-openapi/generated"

	oapiMiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type apiController struct{}

// OpenAPI で定義された (GET /users) の実装
func (a apiController) GetUser(ctx echo.Context) error {
	// OpenApi で生成された User モデルを使ってレスポンスを返す
	return ctx.JSON(http.StatusOK, &oapi.User{
		Id:   1,
		Name: "Taro Yamada",
	})
}

// OpenAPI で定義された (POST /users) の実装
func (a apiController) PostUser(ctx echo.Context) error {
	// リクエストボディを構造体にバインド
	user := &oapi.User{}
	ctx.Bind(&user)
	fmt.Println(user)
	// 200 ステータスのみ返す
	return ctx.NoContent(http.StatusOK)
}

func main() {
	// Echo のインスタンス作成
	e := echo.New()

	// OpenApi 仕様に沿ったリクエストかバリデーションをするミドルウェアを設定
	swagger, err := oapi.GetSwagger()
	if err != nil {
		panic(err)
	}
	e.Use(oapiMiddleware.OapiRequestValidator(swagger))
	// ロガーのミドルウェアを設定
	e.Use(middleware.Logger())
	// APIがエラーで落ちてもリカバーするミドルウェアを設定
	e.Use(middleware.Recover())

	// OpenAPI の仕様を満たす構造体をハンドラーとして登録する
	api := apiController{}
	oapi.RegisterHandlers(e, api)

	// 8080ポートで Echo サーバ起動
	e.Logger.Fatal(e.Start(":8080"))
}
