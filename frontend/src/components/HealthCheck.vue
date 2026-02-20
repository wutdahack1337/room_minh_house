<script setup lang="ts">
import {getHealth} from '../services/api';


const checkHealth = async () => {
    try {
        const status = await getHealth();
        alert(`${status.status}`);
    } catch (error) {
        alert(`${error}`);
    }
};

</script>

<template>
    <div class="health-check">
        <button @click="checkHealth" class="check-btn">
            Check Health
        </button>
    </div>
</template>

<style scoped>
    .check-btn{
        font-size: 16px;
        background-color: #42b946;
        border-radius: 5px;
        cursor: pointer;
    }

</style>