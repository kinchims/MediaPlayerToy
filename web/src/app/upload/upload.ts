import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, ElementRef, OnDestroy, Renderer2, ViewChild } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MPTAPI } from '../../services/api.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-upload',
  imports: [MatButtonModule, CommonModule],
  templateUrl: './upload.html',
  styleUrl: './upload.scss',
})
export class Upload implements AfterViewInit, OnDestroy {
  @ViewChild('upload')
  public fileUpload!: ElementRef<HTMLInputElement>;

  private _subscription: Subscription = new Subscription();

  public constructor(
    private _api: MPTAPI,
    private _renderer: Renderer2,
  ) {}

  public ngAfterViewInit() {
    this._subscription.add(
      this._renderer.listen(this.fileUpload.nativeElement, 'change', (x) => {
        if(this.fileUpload.nativeElement.files && this.fileUpload.nativeElement.files.length > 0){
          this._api.uploadFile(this.fileUpload.nativeElement.files[0]).subscribe()
        }
      }),
    );
  }

  public ngOnDestroy() {
    this._subscription.unsubscribe();
  }

  public showDialog() {
    this.fileUpload.nativeElement.showPicker();
  }
}
