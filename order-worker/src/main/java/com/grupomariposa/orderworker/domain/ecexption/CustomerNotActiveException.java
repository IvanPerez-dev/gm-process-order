package com.grupomariposa.orderworker.ecexption;

public class CustomerNotActiveException extends OrderProcessingException {
    public CustomerNotActiveException(String orderId, String customerId) {
        super(orderId, "Customer " + customerId + " is not active");
    }
}
