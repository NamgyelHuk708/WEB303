package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/database"
	"order-service/models"

	consulapi "github.com/hashicorp/consul/api"
)

func discoverService(serviceName string) (string, error) {
	config := consulapi.DefaultConfig()
	config.Address = "consul:8500"

	consul, err := consulapi.NewClient(config)
	if err != nil {
		return "", err
	}

	services, _, err := consul.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances of %s", serviceName)
	}

	service := services[0].Service
	return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}

type CreateOrderRequest struct {
	UserID uint `json:"user_id"`
	Items  []struct {
		MenuItemID uint `json:"menu_item_id"`
		Quantity   int  `json:"quantity"`
	} `json:"items"`
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Discover user-service via Consul
	userServiceURL, err := discoverService("user-service")
	if err != nil {
		http.Error(w, "User service unavailable", http.StatusServiceUnavailable)
		return
	}

	// Call user-service to validate user exists
	userResp, err := http.Get(fmt.Sprintf("%s/users/%d", userServiceURL, req.UserID))
	if err != nil || userResp.StatusCode != http.StatusOK {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	// Create order
	order := models.Order{
		UserID: req.UserID,
		Status: "pending",
	}

	// Discover menu-service via Consul
	menuServiceURL, err := discoverService("menu-service")
	if err != nil {
		http.Error(w, "Menu service unavailable", http.StatusServiceUnavailable)
		return
	}

	// Validate each menu item by calling menu-service
	for _, item := range req.Items {
		// Get menu item to snapshot price
		menuResp, err := http.Get(fmt.Sprintf("%s/menu/%d", menuServiceURL, item.MenuItemID))
		if err != nil || menuResp.StatusCode != http.StatusOK {
			http.Error(w, "Menu item not found", http.StatusBadRequest)
			return
		}

		var menuItem struct {
			Price float64 `json:"price"`
		}
		json.NewDecoder(menuResp.Body).Decode(&menuItem)

		orderItem := models.OrderItem{
			MenuItemID: item.MenuItemID,
			Quantity:   item.Quantity,
			Price:      menuItem.Price,
		}
		order.OrderItems = append(order.OrderItems, orderItem)
	}

	if err := database.DB.Create(&order).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	var orders []models.Order
	if err := database.DB.Preload("OrderItems").Find(&orders).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
