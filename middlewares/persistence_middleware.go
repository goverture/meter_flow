package middlewares

// import (
// 	"log"
// 	"meter_flow/server"
// 	"net/http"
// )

// func (s *server.Server) withPersistence(next http.HandlerFunc) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         // Execute the main handler
//         next.ServeHTTP(w, r)

//         // Only persist on write operations
//         if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
//             err := s.persistence.Save(&s.resources)
//             if err != nil {
//                 log.Println("Error saving resources:", err)
//             }
//         }
//     }
// }