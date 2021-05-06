<template>
    <b-row class="mb-5">
        <b-col>
            <b-card style="max-width: 30rem; margin: auto;">
                <b-form @submit="onSubmit">
                    <b-form-group id="username-input-group" label="账号" label-for="username-input">
                        <b-form-input id="username-input" v-model="form.username" type="text" placeholder="输入账号" required></b-form-input>
                    </b-form-group>
                    <b-form-group id="password-input-group" label="密码" label-for="password-input">
                        <b-form-input id="password-input" v-model="form.password" type="password" placeholder="输入密码" required></b-form-input>
                    </b-form-group>

                    <b-button type="submit" variant="primary">登录</b-button>
                </b-form>
            </b-card>
        </b-col>
    </b-row>
</template>

<script>
import axios from 'axios';

export default {
        name: 'Login',
        components: {},
        data() {
            return {
                form: {
                    username: '',
                    password: '',
                }
            };
        },
        methods: {
            onSubmit(e) {
                e.preventDefault();
                
                axios.post('/api/auth/login-ldap/', this.form).then(response => {
                   this.ToastSuccess('Login as ' + response.data.name);
                   window.location.href = '/';
                }).catch(error => {
                    this.ErrorBox(error)
                });
            },
        },
        mounted() {
        }
    }
</script>
