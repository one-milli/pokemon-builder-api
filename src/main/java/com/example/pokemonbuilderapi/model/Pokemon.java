package com.example.pokemonbuilderapi.model;

import jakarta.persistence.*;

@Entity
@Table(name = "pokemons")
public class Pokemon {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String name;
    private String type;
    private String baseStats;
}
