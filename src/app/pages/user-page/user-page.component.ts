import { Component } from '@angular/core';
import { MasterClass, Recipe, User } from '../../models/models';
import { MainServiceService } from '../../service/main-service.service';
import { Router, RouterLink } from '@angular/router';

@Component({
  selector: 'app-user-page',
  imports: [RouterLink],
  templateUrl: './user-page.component.html',
  styleUrl: './user-page.component.css'
})
export class UserPageComponent {
  user$: User | null = null;
  masterclasses$: MasterClass[] = [];

  constructor(public mainService: MainServiceService, public router: Router) {

    if (this.mainService.user) {
      this.mainService.loginUser({ email: this.mainService.user?.email, password: this.mainService.user?.password }).subscribe(v => {
        this.user$ = v
      
        this.mainService.getMasterClassesOfUser(v).subscribe(v => this.masterclasses$ = v);
      });
    }
  }

  ngOnInit(): void {
    
  }

  handleMasterclassRegistration(masterclassId: number): void {
    this.mainService.registerForMasterclass(masterclassId, 1).subscribe({
      next: () => console.log('Registration successful'),
      error: (err) => console.error('Registration failed:', err)
    });
  }


  login() {
    if (this.mainService.user) {
      this.mainService.user = null
    } else {
      this.router.navigate(["login"]);
    }
  }

  redirectToProfile(rout: string) {
    this.router.navigate([rout]);
  }
}