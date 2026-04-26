package com.example.orderservice.web.controller;

import com.example.orderservice.application.usecase.ManageCartUseCase;
import com.example.orderservice.domain.model.Cart;
import com.example.orderservice.domain.model.CartItem;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/carts")
@RequiredArgsConstructor
public class CartController {
    private final ManageCartUseCase manageCartUseCase;

    @PostMapping("/{customerId}")
    public Cart addToCart(@PathVariable UUID customerId, @RequestParam UUID restaurantId, @RequestBody CartItem item) {
        return manageCartUseCase.addToCart(customerId, restaurantId, item);
    }

    @GetMapping("/{customerId}")
    public Cart getCart(@PathVariable UUID customerId) {
        return manageCartUseCase.getCart(customerId);
    }
}
