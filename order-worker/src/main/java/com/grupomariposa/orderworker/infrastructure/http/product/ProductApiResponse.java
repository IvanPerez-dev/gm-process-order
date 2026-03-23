package com.grupomariposa.orderworker.infrastructure.http.product;

public record ProductApiResponse(
        String id,
        String name,
        String description,
        Double price
) {}