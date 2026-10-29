package services
import (
 "time"
 "math"
 "fluxo/models"
)
func PriorityScore(t models.Tarefa, now time.Time) float64 {
 points := map[string]float64{"urgente":100,"alta":65,"media":35,"baixa":10}[t.Prioridade]
 if t.Origem=="Ouvidoria" { points+=15 } else { points+=5 }
 points+=math.Min(30,math.Max(0,now.Sub(t.AtualizadoEm).Hours()/24)*3)
 if IsLate(t,now){points+=25}
 due:=Deadline(t);if t.Status!="concluida"&&!due.Before(now)&&due.Sub(now)<=24*time.Hour{points+=10}
 return points
}
