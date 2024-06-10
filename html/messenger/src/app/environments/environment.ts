const environment_dev: environment = {
  base_url: 'https://localhost'
}
const environment_prod: environment = {
  base_url: ''
}

let current_environment: environment = window.location.href.startsWith("http://localhost:4200") ? environment_dev : environment_prod

export function resolve(path: string = "") {
  return current_environment.base_url + "/api" + (path ? `/${path}` : "")
}

interface environment {
  base_url: string
}
