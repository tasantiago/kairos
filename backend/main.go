package main
import (
 "context"
 "errors"
 "log/slog"
 "net/http"
 "os"
 "os/signal"
 "strings"
 "syscall"
 "time"
 "net/mail"
 "fluxo/handlers"
 "fluxo/migrations"
 "fluxo/repositories"
 "fluxo/utils"
 "github.com/gorilla/mux"
 "github.com/jackc/pgx/v5/pgxpool"
 "github.com/joho/godotenv"
 "golang.org/x/crypto/bcrypt"
)
func main(){godotenv.Load();if err:=run();err!=nil{slog.Error("startup failed","error",err);os.Exit(1)}}
func run()error{ctx,stop:=signal.NotifyContext(context.Background(),os.Interrupt,syscall.SIGTERM);defer stop();secret:=os.Getenv("JWT_SECRET");if len(secret)<32||strings.Contains(secret,"replace-with"){return errors.New("JWT_SECRET must be set to at least 32 random characters")};dsn:=os.Getenv("DATABASE_URL");if dsn==""{return errors.New("DATABASE_URL is required")};dbCtx,cancel:=context.WithTimeout(ctx,30*time.Second);defer cancel();db,err:=pgxpool.New(dbCtx,dsn);if err!=nil{return err};defer db.Close();if err=db.Ping(dbCtx);err!=nil{return err};if os.Getenv("AUTO_MIGRATE")!="false"{if err=migrations.Apply(dbCtx,db);err!=nil{return err}};if err=bootstrap(dbCtx,db);err!=nil{return err};upload:=os.Getenv("UPLOAD_DIR");if upload==""{upload="./uploads"};a:=&handlers.API{Repo:&repositories.Repository{DB:db},Secret:[]byte(secret),UploadDir:upload};router:=mux.NewRouter();router.HandleFunc("/health",func(w http.ResponseWriter,r *http.Request){healthCtx,cancel:=context.WithTimeout(r.Context(),2*time.Second);defer cancel();if err:=db.Ping(healthCtx);err!=nil{utils.Error(w,503,"Banco de dados indisponível.");return};utils.JSON(w,200,map[string]string{"status":"ok"})}).Methods("GET");a.Register(router);router.NotFoundHandler=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){utils.Error(w,404,"Endpoint não encontrado.")});port:=os.Getenv("PORT");if port==""{port="8080"};server:=&http.Server{Addr:":"+port,Handler:middleware(router),ReadHeaderTimeout:5*time.Second,ReadTimeout:30*time.Second,WriteTimeout:30*time.Second,IdleTimeout:60*time.Second,MaxHeaderBytes:1<<20};result:=make(chan error,1);go func(){slog.Info("API listening","port",port);result<-server.ListenAndServe()}();select{case err=<-result:if !errors.Is(err,http.ErrServerClosed){return err};case <-ctx.Done():shutdownCtx,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel();return server.Shutdown(shutdownCtx)};return nil}
func middleware(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){start:=time.Now();defer func(){if recovered:=recover();recovered!=nil{slog.Error("request panic","error",recovered);utils.Error(w,500,"Não foi possível concluir a solicitação.")};slog.Info("request","method",r.Method,"path",r.URL.Path,"duration_ms",time.Since(start).Milliseconds())}();w.Header().Set("X-Content-Type-Options","nosniff");w.Header().Set("Referrer-Policy","same-origin");origin:=r.Header.Get("Origin");if origin!=""{w.Header().Add("Vary","Origin");allowed:=false;for _,v:=range strings.Split(os.Getenv("CORS_ORIGIN"),","){if strings.TrimSpace(v)==origin{allowed=true}};if !allowed{utils.Error(w,403,"Origem não autorizada.");return};w.Header().Set("Access-Control-Allow-Origin",origin);w.Header().Set("Access-Control-Allow-Headers","Authorization,Content-Type");w.Header().Set("Access-Control-Allow-Methods","GET,POST,PUT,DELETE,OPTIONS")};if r.Method=="OPTIONS"{w.WriteHeader(204);return};next.ServeHTTP(w,r)})}
func bootstrap(ctx context.Context,db *pgxpool.Pool)error{tx,err:=db.Begin(ctx);if err!=nil{return err};defer tx.Rollback(ctx);if _,err=tx.Exec(ctx,"SELECT pg_advisory_xact_lock(19482631)");err!=nil{return err};if _,err=tx.Exec(ctx,"INSERT INTO configuracoes(id,pontos_por_tarefa) VALUES(1,10) ON CONFLICT DO NOTHING");err!=nil{return err};var count int;if err=tx.QueryRow(ctx,"SELECT count(*) FROM usuarios").Scan(&count);err!=nil{return err};if count>0{return tx.Commit(ctx)};email:=strings.ToLower(strings.TrimSpace(os.Getenv("ADMIN_EMAIL")));password:=os.Getenv("ADMIN_PASSWORD");parsed,parseErr:=mail.ParseAddress(email);if parseErr!=nil||parsed.Address!=email||len(password)<8||len(password)>72{return errors.New("first startup requires valid ADMIN_EMAIL and ADMIN_PASSWORD (8-72 bytes)")};name:=os.Getenv("ADMIN_NAME");if name==""{name="Gestor"};if len([]rune(name))>120{return errors.New("ADMIN_NAME exceeds 120 characters")};hash,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost);if err!=nil{return err};var sector string;err=tx.QueryRow(ctx,"INSERT INTO setores(nome,capacidade,cor) VALUES('Comunicação Social',20,'#32886c') ON CONFLICT(nome) DO UPDATE SET nome=EXCLUDED.nome RETURNING id::text").Scan(&sector);if err!=nil{return err};if _,err=tx.Exec(ctx,"INSERT INTO usuarios(nome,email,senha_hash,setor_id,perfil) VALUES($1,$2,$3,$4,'gestor')",name,email,string(hash),sector);err!=nil{return err};return tx.Commit(ctx)}
