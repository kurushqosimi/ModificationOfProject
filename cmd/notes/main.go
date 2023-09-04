package main

import (
	"log"
	"main/internal/configs"
	"main/internal/configs/repositories"
	"main/internal/handlers"
	"main/internal/services"
	"net/http"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
func run() error {
	config, err := configs.InitConfigs()
	if err != nil {
		return err
	}
	address := config.ServerSetting.Host + config.ServerSetting.Port
	//sqlDB, err := conn.DB.DB()
	//if err != nil {
	//	fmt.Println("Ошибка при получении экземпляра *sql.DB:", err)
	//	return err
	//}
	//defer func() {
	//	err = sqlDB.Close()
	//	if err != nil {
	//		fmt.Println("Ошибка при разрыве соединения с базой данных:", err)
	//	}
	//}()
	conn, err := repositories.GetConnection(config)
	if err != nil {
		return err
	}
	service := services.NewService(conn)
	handler := handlers.NewHandler(service)
	router := handlers.NewRouter(handler)
	log.Println(config)
	srv := http.Server{
		Addr:    address,
		Handler: router,
	}
	err = srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
