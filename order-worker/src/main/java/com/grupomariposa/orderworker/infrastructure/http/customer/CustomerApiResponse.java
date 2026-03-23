package com.grupomariposa.orderworker.infrastructure.http.customer;

public record CustomerApiResponse(   String id,
                                     String name,
                                     String email,
                                     boolean isActive
) {
}
