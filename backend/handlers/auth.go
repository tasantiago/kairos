package handlers
import (
 "context"
 "net/http"
 "net"
 "strings"
 "sync"
 "time"
 "fluxo/models"
 "fluxo/repositories"
 "fluxo/utils"
 "github.com/golang-jwt/jwt/v5"
 "github.com/gorilla/mux"
 "golang.org/x/crypto/bcrypt"
)
type API struct { Repo *repositories.Repository; Secret []byte; UploadDir string; mu sync.Mutex; attempts map[string]attempt }
type attempt struct { count int; start time.Time }
type userKey struct{}
func User(r *http.Request)models.Usuario{return r.Context().Value(userKey{}).(models.Usuario)}
func(a *API)limitLogin(r *http.Request)bool{host,_,err:=net.SplitHostPort(r.RemoteAddr);if err!=nil{host=r.RemoteAddr};a.mu.Lock();defer a.mu.Unlock();if a.attempts==nil{a.attempts=make(map[string]attempt)};for key,value:=range a.attempts{if time.Since(value.start)>15*time.Minute{delete(a.attempts,key)}};v:=a.attempts[host];if v.start.IsZero(){if len(a.attempts)>=10000{return false};v.start=time.Now()};v.count++;a.attempts[host]=v;return v.count<=20}
func(a *API)Login(w http.ResponseWriter,r *http.Request){if !a.limitLogin(r){utils.Error(w,429,"Muitas tentativas. Aguarde 15 minutos.");return};var body struct{Email string `json:"email"`;Senha string `json:"senha"`};if !utils.Decode(w,r,&body){return};var u models.Usuario;var hash string;err:=a.Repo.DB.QueryRow(r.Context(),"SELECT id::text,nome,email,perfil,setor_id::text,senha_hash FROM usuarios WHERE email=$1",strings.ToLower(strings.TrimSpace(body.Email))).Scan(&u.ID,&u.Nome,&u.Email,&u.Perfil,&u.SetorID,&hash);if err!=nil||bcrypt.CompareHashAndPassword([]byte(hash),[]byte(body.Senha))!=nil{utils.Error(w,401,"E-mail ou senha incorretos.");return};now:=time.Now();token,err:=jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.RegisteredClaims{Subject:u.ID,Issuer:"fluxo",Audience:jwt.ClaimStrings{"fluxo-api"},IssuedAt:jwt.NewNumericDate(now),ExpiresAt:jwt.NewNumericDate(now.Add(8*time.Hour))}).SignedString(a.Secret);if err!=nil{utils.Error(w,500,"Não foi possível entrar.");return};utils.JSON(w,200,map[string]any{"token":token,"usuario":u})}
func(a *API)Auth(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){header:=r.Header.Get("Authorization");if !strings.HasPrefix(header,"Bearer "){utils.Error(w,401,"Entre no workspace para continuar.");return};claims:=&jwt.RegisteredClaims{};parsed,err:=jwt.ParseWithClaims(strings.TrimPrefix(header,"Bearer "),claims,func(t *jwt.Token)(any,error){return a.Secret,nil},jwt.WithValidMethods([]string{"HS256"}),jwt.WithIssuer("fluxo"),jwt.WithAudience("fluxo-api"),jwt.WithExpirationRequired());if err!=nil||!parsed.Valid{utils.Error(w,401,"Sua sessão expirou. Entre novamente.");return};var u models.Usuario;err=a.Repo.DB.QueryRow(r.Context(),"SELECT id::text,nome,email,perfil,setor_id::text FROM usuarios WHERE id=$1",claims.Subject).Scan(&u.ID,&u.Nome,&u.Email,&u.Perfil,&u.SetorID);if err!=nil{utils.Error(w,401,"Sessão inválida.");return};next.ServeHTTP(w,r.WithContext(context.WithValue(r.Context(),userKey{},u)))})}
func manager(w http.ResponseWriter,r *http.Request)bool{if User(r).Perfil!="gestor"{utils.Error(w,403,"Esta ação é exclusiva de gestores.");return false};return true}
func(a *API)Register(r *mux.Router){r.HandleFunc("/api/login",a.Login).Methods("POST");protected:=r.PathPrefix("/api").Subrouter();protected.Use(a.Auth);protected.HandleFunc("/me",func(w http.ResponseWriter,r *http.Request){utils.JSON(w,200,User(r))}).Methods("GET");protected.HandleFunc("/tarefas/prioritaria",a.Prioritaria).Methods("GET");protected.HandleFunc("/tarefas",a.Tasks).Methods("GET","POST");protected.HandleFunc("/tarefas/{id}",a.Task).Methods("GET","PUT","DELETE");protected.HandleFunc("/tarefas/{id}/historico",a.History).Methods("GET");protected.HandleFunc("/tarefas/{id}/chat",a.Chat).Methods("GET","POST");protected.HandleFunc("/tarefas/{id}/anexos",a.Attachments).Methods("GET","POST");protected.HandleFunc("/anexos/{id}",a.Attachment).Methods("GET","DELETE");protected.HandleFunc("/setores",a.Sectors).Methods("GET","POST");protected.HandleFunc("/setores/{id}",a.Sectors).Methods("GET","PUT","DELETE");protected.HandleFunc("/usuarios",a.Users).Methods("GET","POST");protected.HandleFunc("/usuarios/{id}",a.Users).Methods("GET","PUT","DELETE");protected.HandleFunc("/config",a.Config).Methods("GET","PUT");protected.HandleFunc("/metricas",a.Metrics).Methods("GET");protected.HandleFunc("/relatorios",a.Report).Methods("GET");protected.HandleFunc("/notificacoes",a.Notifications).Methods("GET")}
