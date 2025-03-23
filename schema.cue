#Schema: {
	name?: string
	age?:  int & <= 150
	usernames?: [...string]
	passwords?: [...string]
  network: {
    cidr: string
    vpn_password: string
  }
}
