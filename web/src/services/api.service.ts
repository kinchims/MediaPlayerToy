import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../environments/environment';
import { IMPTFile } from '../app/reusables/models/MPTFile';
import { BehaviorSubject, Subject } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class MPTAPI {
  public processingStarted: BehaviorSubject<string> = new BehaviorSubject("");
  public processingCompleted: BehaviorSubject<string> = new BehaviorSubject("");

  public constructor(private _httpClient: HttpClient) {
    const eventSource = new EventSource(`http://localhost:7000/sse`);

    eventSource.addEventListener("completed", (x) => {
      this.processingCompleted.next(x.data);
    })

     eventSource.addEventListener("started", (x) => {
      this.processingStarted.next(x.data);
    })
  }

  public getFiles() {
    return this._httpClient.get<IMPTFile[]>(`${environment.api_url}/files`);
  }

  public getFile(file: IMPTFile) {
    return this._httpClient.get(`${environment.api_url}/files/${file.id}`);
  }

  public deleteFile(file: IMPTFile) {
    return this._httpClient.delete(`${environment.api_url}/files/${file.id}`);
  }

  public assignCard(file: IMPTFile) {
    return this._httpClient.post(`${environment.api_url}/files/${file.id}`, null);
  }

  public uploadFile(file: File) {
    let formdata = new FormData();
    formdata.append('file', file);

    return this._httpClient.post(`${environment.api_url}/files`, formdata);
  }
}
