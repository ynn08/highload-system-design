package com.example.orderservice.infrastructure.persistence;

import com.example.orderservice.domain.model.Cart;
import com.example.orderservice.domain.repository.CartRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Component;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.TimeUnit;

@Component
@RequiredArgsConstructor
public class RedisCartRepository implements CartRepository {
    private final RedisTemplate<String, Object> redisTemplate;
    private static final String KEY_PREFIX = "cart:";

    @Override
    public void save(Cart cart) {
        String key = KEY_PREFIX + cart.getCustomerId().toString();
        redisTemplate.opsForValue().set(key, cart, 7, TimeUnit.DAYS);
    }

    @Override
    public Optional<Cart> findByCustomerId(UUID customerId) {
        String key = KEY_PREFIX + customerId.toString();
        Cart cart = (Cart) redisTemplate.opsForValue().get(key);
        return Optional.ofNullable(cart);
    }

    @Override
    public void deleteByCustomerId(UUID customerId) {
        String key = KEY_PREFIX + customerId.toString();
        redisTemplate.delete(key);
    }
}
