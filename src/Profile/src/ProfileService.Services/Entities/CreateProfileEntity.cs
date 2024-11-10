﻿namespace ProfileService.Services.Entities;

public class CreateProfileEntity
{
    public bool Sex { get; set; }
    public int Age { get; set; }
    public string Name { get; set; } = null!;
    public string Description { get; set; } = null!;
    public List<string> Photos { get; set; } = new();
}