package com.grupomariposa.orderworker.infrastructure.mongodb.document;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class OrderItemDocument {
    private String productId;
    private String name;
    private String description;
    private Double price;
}
