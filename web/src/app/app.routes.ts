import { Routes } from '@angular/router';
import { Home } from './home/home';

export const routes: Routes = [
    { path: '', component: Home, children: [
        { path: 'files', loadComponent:  () => import("./files/files").then((x) => x.Files), title: 'Files'},
        { path: 'upload', loadComponent: () => import("./upload/upload").then((x) => x.Upload), title: 'Upload'}
    ] }
];
