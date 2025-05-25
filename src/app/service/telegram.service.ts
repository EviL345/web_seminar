import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TelegramService {

  private apiUrl = 'http://localhost:8080/api/tg'; // Update the port if needed

  constructor(private http: HttpClient) {}

  sendMessage(message: string): Observable<string> {
    const encodedMessage = encodeURIComponent(message); // Ensure safe URL formatting
    return this.http.post(this.apiUrl + '/' + encodedMessage, null, { responseType: 'text' });
  }

  
}
