package com.grupomariposa.orderworker.domain.model;

import lombok.AccessLevel;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;

@Getter
@Builder
@AllArgsConstructor(access = AccessLevel.PRIVATE)
public final class Customer {
    private String id;
    private String name;
    private String email;
    private Boolean isActive;


  public static Customer create(String id, String name, String email, Boolean isActive){

      return new Customer(id,name, email, isActive);
  }

}
