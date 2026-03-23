package com.grupomariposa.orderworker.ecexption;

public class OrderProcessingException extends RuntimeException  {
    private final String orderId;

    public OrderProcessingException(String orderId, String message) {
        super(message);
        this.orderId = orderId;
    }

    public OrderProcessingException(String orderId, String message, Throwable cause) {
        super(message, cause);
        this.orderId = orderId;
    }

    public String getOrderId() { return orderId; }
}
