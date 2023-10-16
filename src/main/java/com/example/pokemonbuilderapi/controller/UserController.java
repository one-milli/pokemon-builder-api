package com.example.pokemonbuilderapi.controller;

import com.example.pokemonbuilderapi.model.User;
import com.example.pokemonbuilderapi.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/users")
public class UserController {

    @Autowired
    private UserService userService;

    @GetMapping("/{username}")
    public User getUserByUsername(@PathVariable String username){
        return userService.findByUername(username);
    }
}
