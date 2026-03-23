package com.grupomariposa.orderworker.ecexption;

public class ProductNotFoundException extends OrderProcessingException {
    public ProductNotFoundException(String orderId, String productId) {
        super(orderId, "Product " + productId + " not found in catalog");
    }
}