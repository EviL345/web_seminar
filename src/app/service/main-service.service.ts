import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Chef, Recipe, MasterClass, User, Registration, RecipeIngredient } from '../models/models';

@Injectable({
  providedIn: 'root'
})
export class MainServiceService {
  private baseUrl = 'http://localhost:8080/api';

  user: User | null = null;

  constructor(private http: HttpClient) {}

  getChefs(): Observable<Chef[]> {
    return this.http.get<Chef[]>(`${this.baseUrl}/chefs`);
  }

  getRecipes(): Observable<Recipe[]> {
    return this.http.get<Recipe[]>(`${this.baseUrl}/recipes`);
  }

  getRecipeById(id: number): Observable<Recipe> {
    return this.http.get<Recipe>(`${this.baseUrl}/recipes/${id}`);
  }

  createRecipe(recipe: Recipe): Observable<Recipe> {
    return this.http.post<Recipe>(`${this.baseUrl}/recipes`, recipe);
  }

  updateRecipe(id: number, recipe: Partial<Recipe>): Observable<Recipe> {
    return this.http.put<Recipe>(`${this.baseUrl}/recipes/${id}`, recipe);
  }

  deleteRecipe(id: number): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/recipes/${id}`);
  }

  getMasterclasses(): Observable<MasterClass[]> {
    return this.http.get<MasterClass[]>(`${this.baseUrl}/masterclasses`);
  }

  registerForMasterclass(masterclassId: number, userId: number): Observable<Registration> {
    return this.http.post<Registration>(`${this.baseUrl}/masterclasses/${masterclassId}/register`, { user_id: userId });
  }

  generateShoppingList(recipeId: number): Observable<RecipeIngredient[]> {
    return this.http.post<RecipeIngredient[]>(`${this.baseUrl}/shopping-list`, { recipe_id: recipeId });
  }

  getRecommendations(userId: number): Observable<MasterClass[]> {
    return this.http.get<MasterClass[]>(`${this.baseUrl}/recommendations/${userId}`);
  }

  getMasterClassesOfUser(user: User): Observable<MasterClass[]> {
    return this.http.post<MasterClass[]>(`${this.baseUrl}/users/getAllMasterClasses`, user);
  }

  registerUser(user: User): Observable<User> {
    return this.http.post<User>(`${this.baseUrl}/users`, user);
  }

  loginUser(user: Partial<User>): Observable<User> {
    return this.http.post<User>(`${this.baseUrl}/users/login`, user);
  }
}