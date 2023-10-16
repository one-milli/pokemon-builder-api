package com.example.pokemonbuilderapi.service;

import com.example.pokemonbuilderapi.model.User;
import com.example.pokemonbuilderapi.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class UserService {

    @Autowired
    private UserRepository userRepository;
    public User findByUername(String username){
        return userRepository.findByUsername(username);
    }
}
