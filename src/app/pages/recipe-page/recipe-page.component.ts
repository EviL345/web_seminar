import { Component } from '@angular/core';
import { Recipe } from '../../models/models';
import { MainServiceService } from '../../service/main-service.service';
import { Router, RouterLink } from '@angular/router';

@Component({
  selector: 'app-recipe-page',
  imports: [RouterLink],
  templateUrl: './recipe-page.component.html',
  styleUrl: './recipe-page.component.css'
})
export class RecipePageComponent {
  recipes: Recipe[] = [];
  loading = true;
  error = '';

  constructor(private mainService: MainServiceService, public router: Router) {}

  ngOnInit(): void {
    this.loadRecipes();
  }

  loadRecipes(): void {
    this.mainService.getRecipes().subscribe({
      next: (data) => {
        this.recipes = data;
        this.loading = false;
      },
      error: (err) => {
        this.error = 'Ошибка загрузки рецептов';
        this.loading = false;
      }
    });
  }

  redirectToProfile(route: string) {
    this.router.navigate([route]);
  }
}