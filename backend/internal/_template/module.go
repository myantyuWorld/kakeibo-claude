package template

import (
	"backend/pkg/transaction"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// Module: パッケージの公開エントリポイント
// - 外部に公開する Resolver / Creator を保持する
// - どちらか片方のみ、または公開するものがないケースもある
type Module struct {
	Resolver *Resolver
	Creator  *Creator
	handler  *handler
}

// NewModule: DI を行い Module を構築する
// - db と必要な resolver 等を外部から受け取る
// - パッケージ内の handler / service / repository の組み立てはここに閉じる
func NewModule(db *sqlx.DB, tx transaction.Transactor, exampleResolver ExampleResolver) *Module {
	repo := newRepository(db)
	svc := newService(repo, tx)
	svc.resolver = exampleResolver
	h := newHandler(svc)

	return &Module{
		Resolver: NewResolver(repo),
		Creator:  NewCreator(svc),
		handler:  h,
	}
}

// RegisterRoutes: Echo のルーティングを登録する
// - echo.Group を受け取り、パッケージ内のハンドラを登録する
func (m *Module) RegisterRoutes(g *echo.Group) {
	g.POST("", m.handler.create)
	g.GET("/:id", m.handler.getByID)
}
