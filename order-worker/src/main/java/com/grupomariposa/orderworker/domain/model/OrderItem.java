package com.grupomariposa.orderworker.domain.model;

import lombok.AllArgsConstructor;
import lombok.Getter;

@Getter
@AllArgsConstructor
public class OrderItem {
    private final String productId;
    private String name;
    private String description;
    private Double price;
    private Integer quantity;


    public Double getTotal() {
        return price != null && quantity != null ? price * quantity : 0.0;
    }
    public static OrderItem of(String productId) {
        return new OrderItem(productId, null, null, null,null);
    }

    public void enrich(String name, String description, Double price) {
        this.name = name;
        this.description = description;
        this.price = price;
    }

    public void setQuantity(Integer quantity){
        this.quantity = quantity;
    }

    public boolean isEnriched() {
        return name != null && price != null;
    }
}
