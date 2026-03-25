import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatTooltipModule } from '@angular/material/tooltip';
import { ActivatedRoute, RouterModule } from '@angular/router';

@Component({
  selector: 'app-nav',
  imports: [MatButtonModule, MatIconModule, CommonModule, RouterModule, MatTooltipModule, MatMenuModule],
  templateUrl: './nav.component.html',
  styleUrl: './nav.component.scss',
})
export class NavComponent {
  @Input()
  public items: INavItem[] = [];

  public handle(item: INavItem) {
    if (item.function) item.function();
  }

  public constructor(public active: ActivatedRoute) {}

  public getClasses(item: INavItem) {
    let classes = { disabled: item.disabled };

    if (item.classes) {
      item.classes.forEach((x) => {
        (classes as any)[x] = true;
      });
    }

    return classes;
  }
}

export interface INavItem {
  path?: string;
  function?: () => void;
  name: string;
  icon?: string;
  classes?: string[];
  disabled?: boolean;
}
