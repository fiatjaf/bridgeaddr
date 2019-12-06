/** @format */

import App from './App.html'
import Donate from './Donate.html'

const main = document.querySelector('main')
const Component = {donate: Donate, app: App}[main.dataset.component]

const app = new Component({
  target: main
})

export default app
