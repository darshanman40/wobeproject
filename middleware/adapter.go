package middleware

//TODO CLEAN UP
//
// //Adapter ....
// type Adapter func(w http.ResponseWriter, r *http.Request) // error
//
// //Adapters ...
// func Adapters(handlers ...Adapter) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		for _, h := range handlers {
// 			h(w, r)
// 		}
// 	})
// }
