import { Component } from '@angular/core';
import { TelegramService } from '../../service/telegram.service';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';

@Component({
  selector: 'app-telegram',
  imports: [FormsModule, RouterLink],
  templateUrl: './telegram.component.html',
  styleUrl: './telegram.component.css'
})
export class TelegramComponent {
  message: string = '';
  chatLog: { sender: 'user' | 'bot', text: string }[] = [];
  isLoading: boolean = false;

  constructor(private telegramBotService: TelegramService, public router: Router) {}

  sendMessage() {
    if (!this.message.trim()) return;

    const userMessage = this.message;
    this.chatLog.push({ sender: 'user', text: userMessage });
    this.message = '';
    this.isLoading = true;

    this.telegramBotService.sendMessage(userMessage).subscribe({
      next: (response) => {
        this.chatLog.push({ sender: 'bot', text: response });
        this.isLoading = false;
      },
      error: () => {
        this.chatLog.push({ sender: 'bot', text: 'Error: Failed to send message.' });
        this.isLoading = false;
      }
    });
  }

  redirectToProfile(route: string) {
    this.router.navigate([route]);
  }
}
