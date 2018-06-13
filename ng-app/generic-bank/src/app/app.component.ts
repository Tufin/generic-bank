import {Component, ElementRef, OnInit, Renderer2, ViewChild} from '@angular/core';
import {ActivatedRoute, ChildActivationEnd, Router} from "@angular/router";
import {DomSanitizer} from "@angular/platform-browser";
import {BehaviorSubject} from "rxjs/internal/BehaviorSubject";
import animations from "./animations/animations";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.less'],
  animations: [ animations.routerTransition, animations.opacityTransition ]
})
export class AppComponent implements OnInit{
  private height$ = new BehaviorSubject(0);
  public heightObs = this.height$.asObservable();
  public height = 0;

  @ViewChild('main') innerView: ElementRef;

  title = 'app';
  showBack = false;
  showTime = false;

  constructor(private activatedRoute: ActivatedRoute,
              private router: Router,
              public safe: DomSanitizer,
              private renderer: Renderer2) {
  }

  ngOnInit() {
    this.router.events.subscribe((event) => {
      // console.log(event);
      if(event instanceof ChildActivationEnd) {
        const data = event.snapshot.firstChild['data'];
        this.title = data['title'];
        this.showBack = data['showBack'];
        this.showTime = data['showTime'] || false;
      }
    });

    this.heightObs.subscribe((height) => {
      this.height = height;
      this.innerView.nativeElement.height = height;
      this.renderer.setStyle(this.innerView.nativeElement, 'height', `${height}px`);
      // console.log(height);
    });
  }

  applyHeight(height) {
    this.height$.next(height);
  }
}
