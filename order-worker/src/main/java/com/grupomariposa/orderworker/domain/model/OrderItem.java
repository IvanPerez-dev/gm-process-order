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


    public static OrderItem of(String productId) {
        return new OrderItem(productId, null, null, null);
    }

    public void enrich(String name, String description, Double price) {
        this.name = name;
        this.description = description;
        this.price = price;
    }

    public boolean isEnriched() {
        return name != null && price != null;
    }
}
