import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, ElementRef, OnDestroy, Renderer2, Signal, signal, ViewChild, WritableSignal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MPTAPI } from '../../services/api.service';
import { Subscription } from 'rxjs';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatProgressBarModule } from '@angular/material/progress-bar';

@Component({
  selector: 'app-upload',
  imports: [MatButtonModule, CommonModule, MatProgressBarModule],
  templateUrl: './upload.html',
  styleUrl: './upload.scss',
})
export class Upload implements AfterViewInit, OnDestroy {
  @ViewChild('upload')
  public fileUpload!: ElementRef<HTMLInputElement>;

  public uploading: WritableSignal<boolean> = signal(false);
  private _subscription: Subscription = new Subscription();

  public constructor(
    private _api: MPTAPI, 
    private _renderer: Renderer2,
    private _snackbar: MatSnackBar
  ) {}

  public ngAfterViewInit() {
    this._subscription.add(
      this._renderer.listen(this.fileUpload.nativeElement, 'change', (x) => {
        if(this.fileUpload.nativeElement.files && this.fileUpload.nativeElement.files.length > 0){
          this.uploading.set(true);
          this._api.uploadFile(this.fileUpload.nativeElement.files[0]).subscribe({
            complete: () => { 
             this.uploading.set(false)
              this._snackbar.open("Successfully uploaded file", "Dismiss", { duration: 1000, horizontalPosition: 'right', verticalPosition: 'top'})
            },
            error: () => { 
              this.uploading.set(false)
              this._snackbar.open("Failed to upload file", "Dismiss", { duration: 1000, horizontalPosition: 'right', verticalPosition: 'top'})
             }
          })
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
