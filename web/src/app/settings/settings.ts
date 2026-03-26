import { Component, signal, WritableSignal } from '@angular/core';
import {
  FormControl,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { MatSliderModule } from '@angular/material/slider';
import { MPTAPI } from '../../services/api.service';
import { ActivatedRoute } from '@angular/router';
import { Config } from '../models/config.model';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-settings',
  imports: [
    MatSliderModule,
    MatFormFieldModule,
    MatInputModule,
    ReactiveFormsModule,
    FormsModule,
    MatIconModule,
    MatButtonModule,
  ],
  templateUrl: './settings.html',
  styleUrl: './settings.scss',
})
export class Settings {
  public form!: WritableSignal<FormGroup>;
  public type: WritableSignal<string> = signal('password');

  public constructor(
    private _api: MPTAPI,
    private _activatedRoute: ActivatedRoute,
  ) {
    this._activatedRoute.data.subscribe((x) => {
      const c = x['config'] as Config;
      this.form = signal(
        new FormGroup({
          volume: new FormControl(c.volume),
          ssid: new FormControl(c.wifiAPName),
          password: new FormControl(null, { validators: [Validators.minLength(8)] }),
        }),
      );
    });
  }

  public formatLabel(value: number): string {
    if (value >= 1000) {
      return value + '%';
    }
    return `${value}`;
  }

  public saveConfig() {
    let form = this.form()
    if (!form.valid || !form.dirty) {
      return;
    }

    let c = {
      volume: form.controls['volume'].value,
      wifiAPName: form.controls['ssid'].value,
      wifiPassword: form.controls['password'].value,
    } as Config;

    this._api.setSettings(c).subscribe((x) => {
      form = new FormGroup({
        volume: new FormControl(c.volume),
        ssid: new FormControl(c.wifiAPName),
        password: new FormControl(null, { validators: [Validators.minLength(8)] }),
      });
    });
  }
}
