import { Component, signal, Signal, WritableSignal } from '@angular/core';
import { INavItem, NavComponent } from '../reusables/nav/nav.component';
import { ActivatedRoute, NavigationEnd, Route, Router, RouterOutlet } from '@angular/router';
import { filter, map, Observable, switchMap } from 'rxjs';
import { routes } from '../app.routes';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-home',
  imports: [NavComponent, RouterOutlet, CommonModule],
  templateUrl: './home.html',
  styleUrl: './home.scss',
})
export class Home {
  public items: INavItem[] = [
    {
      path: 'files',
      icon: 'library_music',
      name: 'files',
    },
    { 
      path: 'upload',
      icon: 'upload',
      name: 'Upload'
    },
     { 
      path: 'settings',
      icon: 'settings',
      name: 'Settings'
    }
  ];
  public title$: WritableSignal<string | undefined> = signal(undefined);
  public constructor(
    private _router: Router,
    private _activatedRoute: ActivatedRoute,
  ) {
    this._router.events
      .pipe(
        filter((event) => event instanceof NavigationEnd),
        map(() => this._activatedRoute),
        map((route) => {
          while (route.firstChild) {
            route = route.firstChild;
          }
          return route;
        }),
        filter((route) => route.outlet === 'primary'),
        switchMap((route) => route.title),
      )
      .subscribe((x) => {
        this.title$.set(x);
      });
  }
}
