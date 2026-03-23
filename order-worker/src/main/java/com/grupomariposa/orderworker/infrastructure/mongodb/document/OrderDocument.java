package com.grupomariposa.orderworker.infrastructure.mongodb.document;

import lombok.Builder;
import lombok.Data;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.time.Instant;
import java.util.List;

@Data
@Builder
@Document(collection = "orders")
public class OrderDocument {
    @Id
    private String id;
    private CustomerDocument customer;
    private List<OrderItemDocument> products;
    private String status;
    private int retryCount;
    private String errorMsg;
    private Instant createdAt;
    private Instant processedAt;
}
