package com.grupomariposa.orderworker.infrastructure.mongodb.document;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class CustomerDocument {
    private String id;
    private String name;
    private String email;
}
