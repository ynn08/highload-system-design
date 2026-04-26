package com.example.orderservice.domain.repository;

import com.example.orderservice.domain.model.Cart;
import java.util.Optional;
import java.util.UUID;

public interface CartRepository {
    void save(Cart cart);
    Optional<Cart> findByCustomerId(UUID customerId);
    void deleteByCustomerId(UUID customerId);
}
