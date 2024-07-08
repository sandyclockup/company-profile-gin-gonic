/**
 * This file is part of the Sandy Andryanto Company Profile Website.
 *
 * @author     Sandy Andryanto <sandy.andryanto.dev@gmail.com>
 * @copyright  2024
 *
 * For the full copyright and license information,
 * please view the LICENSE.md file that was distributed
 * with this source code.
 */

import { Component, OnInit } from '@angular/core';
import { environment } from '../../../environments/environment';
import { Title } from "@angular/platform-browser";
import { PageService } from "../../services/page.service"
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-faq',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './faq.component.html',
  styles: ``
})
export class FaqComponent implements OnInit {

  title = environment.title;
  loading:boolean = true;
  content:any

  constructor(private pageService: PageService, private titleService:Title, private router: Router) {
    this.titleService.setTitle("FAQ | " + this.title);
  }

  ngOnInit(): void {
    this.pageService.ping().subscribe((response: any) => {
        this.loadContent()
    }, (error) => {
        console.log(error)
        this.router.navigate(['/unavailable']);
    });
  }

  loadContent(): void{
    this.pageService.faq().subscribe((response: any) => {
        setTimeout(() => {
            this.content = response;
            this.loading = false;
        }, 1500)
    }, (error) => {
        console.log(error)
        this.router.navigate(['/unavailable']);
    });
  }

}
