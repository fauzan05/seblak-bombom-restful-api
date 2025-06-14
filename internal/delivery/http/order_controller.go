package http

import (
	"fmt"
	"html/template"
	"net/url"
	"seblak-bombom-restful-api/internal/delivery/middleware"
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/usecase"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type OrderController struct {
	Log            *logrus.Logger
	UseCase        *usecase.OrderUseCase
	FrontEndConfig *model.FrontEndConfig
}

func NewOrderController(useCase *usecase.OrderUseCase, logger *logrus.Logger, frontEndConfig *model.FrontEndConfig) *OrderController {
	return &OrderController{
		Log:            logger,
		UseCase:        useCase,
		FrontEndConfig: frontEndConfig,
	}
}

func (c *OrderController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateOrderRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}
	getLang := ctx.Query("lang", string(enum_state.ENGLISH))
	request.Lang = enum_state.Languange(getLang)
	getTimeZoneUser := ctx.Query("timezone", "UTC")
	loc, err := time.LoadLocation(getTimeZoneUser)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.TimeZone = *loc
	auth := middleware.GetCurrentUser(ctx)
	request.UserId = auth.ID
	request.FirstName = auth.FirstName
	request.LastName = auth.LastName
	request.Email = auth.Email
	request.Phone = auth.Phone
	if auth.Addresses == nil {
		c.Log.Warnf("address not found/selected!")
		return fiber.NewError(fiber.StatusBadRequest, "address not found/selected!")
	}

	for _, address := range auth.Addresses {
		if address.IsMain {
			request.CompleteAddress = address.CompleteAddress
			if address.Delivery.ID == 0 && request.IsDelivery {
				c.Log.Warnf("delivery not found, please selected one!")
				return fiber.NewError(fiber.StatusNotFound, "delivery not found, please selected one!")
			}
			request.DeliveryId = address.Delivery.ID
		}
	}

	request.CurrentBalance = auth.Wallet.Balance
	request.BaseFrontEndURL = c.FrontEndConfig.BaseURL
	response, err := c.UseCase.Add(ctx, request)
	if err != nil {
		c.Log.Warnf("failed to create a new order : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.ApiResponse[*model.OrderResponse]{
		Code:   201,
		Status: "success to create a new order",
		Data:   response,
	})
}

func (c *OrderController) GetAllCurrent(ctx *fiber.Ctx) error {
	auth := middleware.GetCurrentUser(ctx)
	orderRequest := new(model.GetOrderByCurrentRequest)
	orderRequest.UserId = auth.ID
	response, err := c.UseCase.GetAllCurrent(ctx.Context(), orderRequest)
	if err != nil {
		c.Log.Warnf("failed to get all orders by current user : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*[]model.OrderResponse]{
		Code:   200,
		Status: "success to get all orders by current user",
		Data:   response,
	})
}

func (c *OrderController) GetOrderById(ctx *fiber.Ctx) error {
	getId := ctx.Params("orderId")
	orderId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("failed to convert order_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to convert order_id to integer : %+v", err))
	}
	auth := middleware.GetCurrentUser(ctx)
	
	response, err := c.UseCase.GetOrderById(ctx.Context(), uint64(orderId), auth)
	if err != nil {
		c.Log.Warnf("failed to get order by id : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.OrderResponse]{
		Code:   200,
		Status: "success to get order by id",
		Data:   response,
	})
}

func (c *OrderController) GetAll(ctx *fiber.Ctx) error {
	auth := middleware.GetCurrentUser(ctx)

	search := ctx.Query("search", "")
	trimSearch := strings.TrimSpace(search)

	// ambil data sorting
	getColumn := ctx.Query("column", "")
	getSortBy := ctx.Query("sort_by", "desc")

	// Ambil query parameter 'per_page' dengan default value 10 jika tidak disediakan
	perPage, err := strconv.Atoi(ctx.Query("per_page", "10"))
	if err != nil {
		c.Log.Warnf("invalid 'per_page' parameter : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid 'per_page' parameter : %+v", err))
	}

	// Ambil query parameter 'page' dengan default value 1 jika tidak disediakan
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		c.Log.Warnf("invalid 'page' parameter : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid 'page' parameter : %+v", err))
	}

	response, totalOrders, totalPages, err := c.UseCase.GetAllPaginate(ctx.Context(), page, perPage, trimSearch, getColumn, getSortBy, auth)
	if err != nil {
		c.Log.Warnf("failed to get all orders by current user : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponsePagination[*[]model.OrderResponse]{
		Code:         200,
		Status:       "success to get all orders by current user",
		Data:         response,
		TotalDatas:   totalOrders,
		TotalPages:   totalPages,
		CurrentPages: page,
		DataPerPages: perPage,
	})
}

func (c *OrderController) UpdateOrderStatus(ctx *fiber.Ctx) error {
	getId := ctx.Params("orderId")
	orderId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("failed to convert order_id to integer : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to convert order_id to integer : %+v", err))
	}

	request := new(model.UpdateOrderRequest)
	request.ID = uint64(orderId)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("cannot parse data : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("cannot parse data : %+v", err))
	}
	getLang := ctx.Query("lang", string(enum_state.ENGLISH))
	request.Lang = enum_state.Languange(getLang)
	getTimeZoneUser := ctx.Query("timezone", "UTC")
	loc, err := time.LoadLocation(getTimeZoneUser)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.TimeZone = *loc
	auth := middleware.GetCurrentUser(ctx)
	request.BaseFrontEndURL = c.FrontEndConfig.BaseURL
	response, err := c.UseCase.EditOrderStatus(ctx.Context(), request, auth)
	if err != nil {
		c.Log.Warnf("failed to update order status by selected order : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*model.OrderResponse]{
		Code:   200,
		Status: "success to update order status by selected order",
		Data:   response,
	})
}

func (c *OrderController) GetAllByUserId(ctx *fiber.Ctx) error {
	getId := ctx.Params("userId")
	userId, err := strconv.Atoi(getId)
	if err != nil {
		c.Log.Warnf("failed to convert order id : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to convert order id : %+v", err))
	}

	orderRequest := new(model.GetOrdersByUserIdRequest)
	orderRequest.ID = uint64(userId)
	response, err := c.UseCase.GetByUserId(ctx.Context(), orderRequest)
	if err != nil {
		c.Log.Warnf("failed to get all orders by user id : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.ApiResponse[*[]model.OrderResponse]{
		Code:   200,
		Status: "success to get all orders by user id",
		Data:   response,
	})
}

func (c *OrderController) ShowInvoiceByOrderId(ctx *fiber.Ctx) error {
	getInvoiceId := ctx.Params("invoiceId")
	decoded, err := url.QueryUnescape(getInvoiceId)
	if err != nil {
		c.Log.Warnf("failed to decode invoice id : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to decode invoice id : %+v", err))
	}

	order, app, err := c.UseCase.GetInvoice(ctx.Context(), decoded)
	if err != nil {
		c.Log.Warnf("failed to get order by order id : %+v", err)
		return err
	}

	items := []map[string]any{}
	for _, orderProduct := range order.OrderProducts {
		item := map[string]any{
			"Name":       orderProduct.ProductName,
			"Quantity":   orderProduct.Quantity,
			"UnitPrice":  helper_others.FormatNumberFloat32(orderProduct.Price),
			"TotalPrice": helper_others.FormatNumberFloat32(orderProduct.Price * float32(orderProduct.Quantity)),
		}
		items = append(items, item)
	}

	templatePath := "../internal/templates/pdf/orders/invoice.html"
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	getTimeZoneUser := ctx.Query("timezone", "UTC")
	loc, err := time.LoadLocation(getTimeZoneUser)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	timeZone := helper_others.TimeZoneMap[getTimeZoneUser]
	logoImage := fmt.Sprintf("../uploads/images/application/%s", app.LogoFilename)
	logoImageToBase64, err := helper_others.ImageToBase64(logoImage)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	paymentStatusColor := helper_others.GetPaymentStatusColor(order.PaymentStatus)
	orderMethod := "Ambil Di Tempat (Pickup)"
	if order.IsDelivery {
		orderMethod = "Diantar Ke Alamat"
	}

	bodyBuilder := new(strings.Builder)
	err = tmpl.Execute(bodyBuilder, map[string]any{
		"InvoiceNumber":      order.Invoice,
		"PurchaseDate":       order.CreatedAt.ToTime().In(loc).Format("02 January 2006 15:04"),
		"TimeZone":           timeZone,
		"BuyerName":          order.FirstName + " " + order.LastName,
		"ShippingAddress":    order.CompleteAddress,
		"Items":              items,
		"IsDelivery":         order.IsDelivery,
		"Subtotal":           helper_others.FormatNumberFloat32(order.TotalProductPrice),
		"Discount":           helper_others.FormatNumberFloat32(order.TotalDiscount),
		"ShippingCost":       helper_others.FormatNumberFloat32(order.DeliveryCost),
		"TotalBilling":       helper_others.FormatNumberFloat32(order.TotalFinalPrice + app.ServiceFee),
		"ServiceFee":         helper_others.FormatNumberFloat32(app.ServiceFee),
		"PaymentMethod":      order.PaymentMethod,
		"PaymentStatus":      order.PaymentStatus,
		"PaymentStatusColor": paymentStatusColor,
		"OrderMethod":        orderMethod,
		"UpdatedAt":          order.UpdatedAt.ToTime().In(loc).Format("02 January 2006 15:04"),
		"CompanyTitle":       app.AppName,
		"CompanyPhone":       app.PhoneNumber,
		"CompanyEmail":       app.Email,
		"CompanyAddress":     app.Address,
		"LogoImage":          logoImageToBase64,
	})

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Type("html", "utf-8")
	return ctx.SendString(bodyBuilder.String())
}
