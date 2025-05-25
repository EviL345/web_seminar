import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { User } from '../../models/models';
import { MainServiceService } from '../../service/main-service.service';
import { Router } from '@angular/router';


@Component({
  selector: 'app-login-form',
  imports: [ReactiveFormsModule],
  templateUrl: './login-form.component.html',
  styleUrl: './login-form.component.css'
})
export class LoginComponent implements OnInit {
  loginForm!: FormGroup;
  isLoading = false;
  errorMessage = '';

  constructor(
    private fb: FormBuilder,
    private mainService: MainServiceService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.initLoginForm();
  }

  private initLoginForm(): void {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  onSubmit(): void {
    if (this.loginForm.invalid) {
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';

    this.mainService.loginUser(this.loginForm.value)
      .subscribe({
        next: (response: User) => {
          this.handleSuccessfulLogin(response);
        },
        error: (error: any) => {
          this.handleError(error);
        }
      });
  }

  private handleSuccessfulLogin(user: User): void {
    this.mainService.user = user;
    console.log(user)
    this.router.navigate(["profile"]);
  }

  private handleError(error: any): void {
    this.isLoading = false;
    if (error.error?.message) {
      this.errorMessage = error.error.message;
    } else {
      this.errorMessage = 'An unexpected error occurred';
    }
  }

  get email() { return this.loginForm.get('email'); }
  get password() { return this.loginForm.get('password'); }
}