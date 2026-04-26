package com.example.orderservice.application.usecase;

import com.example.orderservice.domain.model.Cart;
import com.example.orderservice.domain.model.CartItem;
import com.example.orderservice.domain.repository.CartRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class ManageCartUseCase {
    private final CartRepository cartRepository;

    public Cart addToCart(UUID customerId, UUID restaurantId, CartItem item) {
        Cart cart = cartRepository.findByCustomerId(customerId)
                .orElse(Cart.builder().customerId(customerId).restaurantId(restaurantId).build());
        
        if (!cart.getRestaurantId().equals(restaurantId)) {
            // Business rule: one restaurant per cart. Reset if different.
            cart = Cart.builder().customerId(customerId).restaurantId(restaurantId).build();
        }
        
        cart.addItem(item);
        cartRepository.save(cart);
        return cart;
    }

    public Cart getCart(UUID customerId) {
        return cartRepository.findByCustomerId(customerId).orElse(new Cart());
    }
}
