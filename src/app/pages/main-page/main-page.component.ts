import { Component, OnInit } from '@angular/core';
import { MainServiceService } from '../../service/main-service.service';
import { MasterClass } from '../../models/models';
import { RouterLink } from '@angular/router';

@Component({
  imports: [RouterLink],
  selector: 'app-main-page',
  templateUrl: './main-page.component.html',
  styleUrls: ['./main-page.component.css']
})
export class MainPageComponent implements OnInit {
  masterclasses: MasterClass[] = [];

  constructor(private apiService: MainServiceService) {}

  ngOnInit(): void {
    this.loadMasterclasses();
  }

  loadMasterclasses(): void {
    this.apiService.getMasterclasses()
    .subscribe(
      (data) => { this.masterclasses = data; },
      (error) => { console.error('Ошибка загрузки мастер-классов:', error); }
    );
  }

  formatDate(dateString: string): string {
    const date = new Date(dateString);
    const options: Intl.DateTimeFormatOptions = {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        hour12: false
    };
    
    const formatter = new Intl.DateTimeFormat('en-GB', options);
    return formatter.format(date);
}

  trackById(index: number, item: MasterClass): number {
    return item.id;
  }
}
