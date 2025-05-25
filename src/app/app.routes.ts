import { Routes } from '@angular/router';
import { MainPageComponent } from './pages/main-page/main-page.component';
import { RecipePageComponent } from './pages/recipe-page/recipe-page.component';
import { UserPageComponent } from './pages/user-page/user-page.component';
import { LoginComponent } from './pages/login-form/login-form.component';
import { TelegramComponent } from './pages/telegram/telegram.component';

export const routes: Routes = [
    { path: '', component: MainPageComponent },
    { path: 'recipes', component: RecipePageComponent }, 
    { path: 'profile', component: UserPageComponent },
    { path: 'login', component: LoginComponent },
    { path: 'bot', component: TelegramComponent}
];
