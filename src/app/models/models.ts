export interface Chef {
    id: number;
    name: string;
    bio?: string;
    specialization?: string;
  }

  
  export interface Ingredient {
    id: number;
    name: string;
  }

export interface MasterClass {
  id: number;
  title: string;
  dateTime: string;
  chef: Chef;
}


export interface Recipe {
  id: number;
  title: string;
  description?: string;
  chef: Chef;
}



export interface RecipeIngredient {
  recipe: Recipe;
  ingredient: Ingredient;
  quantity: string;
}


export interface Registration {
  user: User;
  masterclass: MasterClass;
}

export interface User {
    id: number;
    name: string;
    email: string;
    password: string;
    preferences?: string;
  }
