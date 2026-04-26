package com.example.orderservice.domain.model;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Cart {
    private UUID customerId;
    private UUID restaurantId;
    private List<CartItem> items = new ArrayList<>();

    public void addItem(CartItem item) {
        if (items == null) items = new ArrayList<>();
        items.add(item);
    }
}
