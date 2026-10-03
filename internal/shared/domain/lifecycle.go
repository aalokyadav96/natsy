package domain

const (
	OrderStatusCreated   = "created"
	OrderStatusValidated = "validated"
	OrderStatusReserved  = "reserved"
	OrderStatusPaid      = "paid"
	OrderStatusFulfilled = "fulfilled"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"

	PurchaseOrderStatusDraft     = "draft"
	PurchaseOrderStatusApproved  = "approved"
	PurchaseOrderStatusOrdered   = "ordered"
	PurchaseOrderStatusReceived  = "received"
	PurchaseOrderStatusMatched   = "matched"
	PurchaseOrderStatusPaid      = "paid"
	AccountsPayableStatusPosted  = "posted"
	AccountsPayableStatusPending = "pending"

	PaymentBatchStatusDraft      = "draft"
	PaymentBatchStatusApproved   = "approved"
	PaymentBatchStatusPaid       = "paid"
	PaymentBatchStatusReconciled = "reconciled"

	PaymentStatusPaid       = "paid"
	PaymentStatusReconciled = "reconciled"
)
